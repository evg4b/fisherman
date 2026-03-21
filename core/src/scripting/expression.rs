use anyhow::Result;
use rhai::{Engine, Scope};
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use std::collections::HashMap;

thread_local! {
    static ENGINE: Engine = Engine::new();
}

#[derive(Debug, Clone)]
pub struct Expression {
    pub condition: String,
}

impl<'de> Deserialize<'de> for Expression {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        Ok(Expression::new(s.as_str()))
    }
}

impl Serialize for Expression {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        self.condition.serialize(serializer)
    }
}

impl Expression {
    pub fn new(condition: &str) -> Self {
        Self {
            condition: condition.to_string(),
        }
    }

    pub fn check(&self, variables: &HashMap<String, String>) -> Result<bool> {
        ENGINE.with(|engine| {
            let mut scope = Scope::with_capacity(variables.len());

            variables.iter().for_each(|(key, value)| {
                scope.push(key.to_owned(), value.to_owned());
            });

            engine
                .eval_expression_with_scope::<bool>(&mut scope, &self.condition)
                .map_err(|err| anyhow::anyhow!("Expression error: {}", err))
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_check_expression() {
        let actual = Expression::new("1 > 0")
            .check(&HashMap::new())
            .unwrap();

        assert!(actual);
    }

    #[test]
    fn test_expression_returns_false_for_undefined_variable() {
        let actual = Expression::new("is_def_var(\"xx\") && xx > 10")
            .check(&HashMap::new())
            .unwrap();

        assert!(!actual);
    }

    #[test]
    fn test_check_expression_error() {
        let actual = Expression::new("1 >").check(&HashMap::new()).unwrap_err();

        assert_eq!(
            actual.to_string(),
            "Expression error: Syntax error: Script is incomplete (line 1, position 4)"
        );
    }

    #[test]
    fn test_expression_with_integer_parsing() {
        let mut variables = HashMap::new();
        variables.insert("xx".to_string(), "20".to_string());
        let actual = Expression::new("parse_int(xx) > 10")
            .check(&variables)
            .unwrap();

        assert!(actual);
    }

    #[test]
    fn test_expression_with_complex_or_condition() {
        let mut variables = HashMap::new();
        variables.insert("xx".to_string(), "91".to_string());
        let actual = Expression::new("(is_def_var(\"yy\") && parse_int(yy) > 10) || (is_def_var(\"xx\") && parse_int(xx) > 10)")
            .check(&variables)
            .unwrap();

        assert!(actual);
    }

    #[test]
    fn expression_deserialize() {
        let rule = r#"
            "is_def_var(\"test\")"
        "#;

        let rule: Expression = serde_json::from_str(rule).unwrap();

        assert_eq!(rule.condition, "is_def_var(\"test\")".to_string());
    }

    #[test]
    fn expression_deserialize_rejects_object_input() {
        let rule = r#"
            { "condition": "1 > 0" }
        "#;

        let err = serde_json::from_str::<Expression>(rule).unwrap_err();

        assert_eq!(
            err.to_string(),
            "invalid type: map, expected a string at line 2 column 12"
        );
    }
}
