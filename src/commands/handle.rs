use crate::common::{CommonError, R};
use crate::configuration::Configuration;
use crate::context::Context;
use crate::err;
use crate::hooks::GitHook;
use crate::rules::{Rule, RuleResult};
use crate::ui::hook_display::hook_display;
use regex::Regex;
use std::collections::HashMap;
use std::process::exit;

pub(crate) fn handle_command(context: &impl Context, hook: &GitHook) -> R {
    let config = Configuration::load(context.repo_path())?;
    println!("{}", hook_display(hook, config.files));

    let variables = extract_variables(context, config.extract)?;

    match config.hooks.get(hook) {
        Some(rules) => {
            let rules_to_exec: Vec<Rule> =
                rules.iter().map(|rule| Rule::new(rule.clone(), variables.clone())).collect();

            let results: Vec<RuleResult> =
                rules_to_exec.into_iter().map(|rule| rule.exec()).collect();

            for rule in results.iter() {
                match rule {
                    RuleResult::Success { name } => {
                        println!("{} executed successfully", name);
                    }
                    RuleResult::Failure { message, name } => {
                        eprintln!("{}: {}", name, message);
                    }
                }
            }

            if results
                .iter()
                .any(|r| matches!(r, RuleResult::Failure { .. }))
            {
                exit(1);
            }
        }
        None => println!("No rules found for hook {}", hook),
    };

    Ok(())
}

#[derive(Eq, Hash, PartialEq)]
enum VariableSource {
    Branch { optional: bool },
    RepoPath { optional: bool },
}

impl VariableSource {
    fn from_str(s: &str) -> R<Self> {
        match s.to_lowercase().as_str() {
            "branch" => Ok(VariableSource::Branch { optional: false }),
            "branch?" => Ok(VariableSource::Branch { optional: true }),
            "repo_path" => Ok(VariableSource::RepoPath { optional: false }),
            "repo_path?" => Ok(VariableSource::RepoPath { optional: true }),
            _ => err!(CommonError::new(format!("Invalid variable source: {}", s))),
        }
    }

    fn extract(&self, context: &impl Context) -> R<(String, bool)> {
        match self {
            VariableSource::Branch { optional } => Ok((context.current_branch()?, *optional)),
            VariableSource::RepoPath { optional } => {
                Ok((context.repo_path().to_string_lossy().to_string(), *optional))
            }
        }
    }
}

fn transform_array(arr: Vec<String>) -> R<HashMap<VariableSource, Vec<Regex>>> {
    let mut map: HashMap<VariableSource, Vec<Regex>> = HashMap::new();

    for entry in arr {
        match entry.split_once(":") {
            Some((key, value)) => {
                let expression = Regex::new(value)?;
                let key = VariableSource::from_str(key)?;
                map.entry(key).or_default().push(expression);
            }
            None => err!(CommonError::new("Invalid extract format")),
        }
    }

    Ok(map)
}

fn extract_variables(
    ctx: &impl Context,
    extract: Vec<String>,
) -> R<HashMap<String, String>> {
    let expressions = transform_array(extract)?;

    let mut variables = HashMap::new();

    for (source, expressions) in expressions.iter() {
        let (source, _optional) = source.extract(ctx)?;
        for expression in expressions.iter() {
            let names = expression.capture_names();
            let captures = expression.captures(&source);
            match captures {
                Some(captures) => {
                    names.into_iter().for_each(|name| {
                        if let Some(name) = name {
                            let demo = captures.name(name);
                            if let Some(demo) = demo {
                                variables.insert(name.to_string(), demo.as_str().to_string());
                            }
                        }
                    });
                }
                None => err!(CommonError::new(
                    format!("The expression {:?} does not match the source {:?}", expression, source)
                )),
            }
        }
    }

    Ok(variables)
}

#[cfg(test)]
mod extract_variables_tests {
    use super::*;
    use crate::context::MockContext;
    use assertor::*;
    use std::path::Path;

    #[test]
    fn accept_empty_vec() {
        let context = &MockContext::new();
        let result = extract_variables(context, vec![]).unwrap();
        assert_that!(result).is_equal_to(HashMap::new());
    }

    #[test]
    fn extract_variables_from_branch() {
        let mut context = MockContext::new();
        context
            .expect_current_branch()
            .returning(|| Ok("master".to_string()));

        let extract = vec!["branch:m(?P<part>.*)".to_string()];

        let result: HashMap<String, String> =
            extract_variables(&context, extract).unwrap();

        assert_that!(result).has_length(1);
        assert_that!(result["part"]).is_equal_to("aster".to_string())
    }

    #[test]
    fn extract_variables_from_path() {
        let mut context = MockContext::new();
        let path = Path::new("/path/to/repo").to_path_buf();
        context.expect_repo_path().return_const(path);

        let extract = vec!["repo_path:^/path/(?P<demo>.*)/repo$".to_string()];

        let result: HashMap<String, String> =
            extract_variables(&context, extract).unwrap();

        assert_that!(result).has_length(1);
        assert_that!(result["demo"]).is_equal_to("to".to_string())
    }

    #[test]
    fn extract_multiple_variables() {
        let mut context = MockContext::new();
        let path = Path::new("/path/to/repo").to_path_buf();
        context.expect_repo_path().return_const(path);

        let extract = vec!["repo_path:^/(?P<S1>.\\S+)/(?P<S2>.\\S+)/(?P<S3>.\\S+)$".to_string()];

        let result: HashMap<String, String> =
            extract_variables(&context, extract).unwrap();

        assert_that!(result).has_length(3);
        assert_that!(result["S1"]).is_equal_to("path".to_string());
        assert_that!(result["S2"]).is_equal_to("to".to_string());
        assert_that!(result["S3"]).is_equal_to("repo".to_string());
    }

    #[test]
    fn should_return_error_when_expression_doesnt_match() {
        let mut context = MockContext::new();
        context
            .expect_current_branch()
            .returning(|| Ok("master".to_string()));

        let extract = vec!["branch:^.&".to_string()];

        let error = extract_variables(&context, extract);

        assert_that!(error).is_err();
        assert_that!(error.unwrap_err().to_string())
            .is_equal_to("The expression Regex(\"^.&\") does not match the source \"master\"".to_string());
    }
}
