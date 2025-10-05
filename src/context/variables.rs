use crate::context::Context;
use anyhow::{bail, Result};
use regex::Regex;
use std::collections::HashMap;

#[derive(Eq, Hash, PartialEq)]
enum VariableSource {
    Branch { optional: bool },
    RepoPath { optional: bool },
}

impl VariableSource {
    fn from_str(s: &str) -> Result<Self> {
        match s.to_lowercase().as_str() {
            "branch" => Ok(VariableSource::Branch { optional: false }),
            "branch?" => Ok(VariableSource::Branch { optional: true }),
            "repo_path" => Ok(VariableSource::RepoPath { optional: false }),
            "repo_path?" => Ok(VariableSource::RepoPath { optional: true }),
            _ => bail!(
                "Unknown source: '{}'. Use 'branch', 'branch?', 'repo_path', or 'repo_path?'",
                s
            ),
        }
    }

    fn extract(&self, context: &impl Context) -> Result<(String, bool)> {
        match self {
            VariableSource::Branch { optional } => Ok((context.current_branch()?, *optional)),
            VariableSource::RepoPath { optional } => {
                Ok((context.repo_path().to_string_lossy().to_string(), *optional))
            }
        }
    }
}

fn transform_array(arr: &[String]) -> Result<HashMap<VariableSource, Vec<Regex>>> {
    let mut map: HashMap<VariableSource, Vec<Regex>> = HashMap::new();

    for entry in arr {
        match entry.split_once(":") {
            Some((key, value)) => {
                let expression = Regex::new(value)?;
                let key = VariableSource::from_str(key)?;
                map.entry(key).or_default().push(expression);
            }
            None => bail!("Invalid format. Use 'source:regex' (e.g., 'branch:^feat/.*')"),
        }
    }

    Ok(map)
}

pub fn extract_variables(
    ctx: &impl Context,
    extract: &[String],
) -> Result<HashMap<String, String>> {
    let expressions = transform_array(extract)?;
    let mut variables = HashMap::with_capacity(expressions.len());

    for (source, expressions) in expressions.iter() {
        let (source, optional) = source.extract(ctx)?;
        for expression in expressions.iter() {
            let names = expression.capture_names();
            let captures = expression.captures(&source);
            match captures {
                Some(captures) => {
                    names.flatten().for_each(|name| {
                        if let Some(demo) = captures.name(name) {
                            variables.insert(name.to_string(), demo.as_str().to_string());
                        }
                    });
                }
                None => {
                    if optional {
                        continue;
                    }

                    bail!("Pattern '{}' doesn't match '{}'", expression, source);
                }
            }
        }
    }

    Ok(variables)
}

#[cfg(test)]
mod extract_variables_tests {
    use super::*;
    use crate::context::MockContext;
    use assert2::assert;
    use std::path::Path;

    #[test]
    fn accept_empty_vec() {
        let context = &MockContext::new();
        let result = extract_variables(context, &[]).unwrap();
        assert!(result.is_empty());
    }

    #[test]
    fn extract_variables_from_branch() {
        let mut context = MockContext::new();
        context
            .expect_current_branch()
            .returning(|| Ok("master".to_string()));

        let extract = vec!["branch:m(?P<part>.*)".to_string()];

        let result: HashMap<String, String> = extract_variables(&context, &extract).unwrap();

        assert!(result.len() == 1);
        assert!(result["part"] == "aster");
    }

    #[test]
    fn extract_variables_from_path() {
        let mut context = MockContext::new();
        let path = Path::new("/path/to/repo").to_path_buf();
        context.expect_repo_path().return_const(path);

        let extract = vec!["repo_path:^/path/(?P<demo>.*)/repo$".to_string()];

        let result: HashMap<String, String> = extract_variables(&context, &extract).unwrap();

        assert!(result.len() == 1);
        assert!(result["demo"] == "to");
    }

    #[test]
    fn extract_multiple_variables() {
        let mut context = MockContext::new();
        let path = Path::new("/path/to/repo").to_path_buf();
        context.expect_repo_path().return_const(path);

        let extract = vec!["repo_path:^/(?P<S1>.\\S+)/(?P<S2>.\\S+)/(?P<S3>.\\S+)$".to_string()];

        let result: HashMap<String, String> = extract_variables(&context, &extract).unwrap();

        assert!(result.len() == 3);
        assert!(result["S1"] == "path");
        assert!(result["S2"] == "to");
        assert!(result["S3"] == "repo");
    }

    #[test]
    fn should_return_error_when_expression_doesnt_match() {
        let mut context = MockContext::new();
        context
            .expect_current_branch()
            .returning(|| Ok("master".to_string()));

        let extract = vec!["branch:^.&".to_string()];

        let error = extract_variables(&context, &extract);

        assert!(error.is_err());
        assert!(error.unwrap_err().to_string() == "Pattern '^.&' doesn't match 'master'");
    }

    #[test]
    fn should_return_not_error_when_expression_is_optional() {
        let mut context = MockContext::new();
        context
            .expect_current_branch()
            .returning(|| Ok("master".to_string()));

        let extract = vec!["branch?:^.&".to_string()];

        let result = extract_variables(&context, &extract).unwrap();

        assert!(result.is_empty());
    }

    #[test]
    fn should_return_not_error_when_eфывфывxpression_is_optional() {
        let mut context = MockContext::new();
        context
            .expect_current_branch()
            .returning(|| Ok("CLIC-48484-bing-front".to_string()));

        let extract = vec!["branch?:^(?<IssueNumber>CLIC-\\d+)-.*$".to_string()];

        let result = extract_variables(&context, &extract).unwrap();

        assert!(result.len() == 1);
        assert!(result["IssueNumber"] == "CLIC-48484");
    }
}
