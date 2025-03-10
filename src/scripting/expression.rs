use anyhow::Result;
use rhai::{Engine, Scope};
use std::collections::HashMap;

pub fn check_expression(expression: &str, variables: &HashMap<String, String>) -> Result<bool> {
    let engine = Engine::new();

    let mut scope = Scope::new();

    variables.iter().for_each(|(key, value)| {
        scope.push(key.clone(), value.clone());
    });

    match engine.eval_expression_with_scope::<bool>(&mut scope, expression) {
        Ok(result) => Ok(result),
        Err(err) => Err(anyhow::anyhow!("Error: {}", err)),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_check_expression() {
        let a = check_expression("1 > 0", &HashMap::new()).unwrap();
        assert_eq!(a, true);
    }

    #[test]
    fn test_check_expression2() {
        let a = check_expression("is_def_var(\"xx\") && xx > 10", &HashMap::new()).unwrap();
        assert_eq!(a, false);
    }

    #[test]
    fn test_check_expression_error() {
        let a = check_expression("1 >", &HashMap::new()).unwrap_err();
        assert_eq!(
            a.to_string(),
            "Error: Syntax error: Script is incomplete (line 1, position 4)"
        );
    }

    #[test]
    fn test_check_expression3() {
        let mut variables = HashMap::new();
        variables.insert("xx".to_string(), "20".to_string());
        let a = check_expression("parse_int(xx) > 10", &variables).unwrap();
        assert_eq!(a, true);
    }

    #[test]
    fn test_check_expression23() {
        let mut variables = HashMap::new();
        variables.insert("xx".to_string(), "91".to_string());
        let a = check_expression("(is_def_var(\"yy\") && parse_int(yy) > 10) || (is_def_var(\"xx\") && parse_int(xx) > 10)", &variables).unwrap();
        assert_eq!(a, true);
    }
}
