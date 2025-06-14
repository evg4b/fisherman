use crate::templates::TemplateString;
use regex::Regex;

pub fn match_expression(expression: &TemplateString, text: &str) -> anyhow::Result<bool> {
    let filled_expression = expression.to_string()?;
    let regex = Regex::new(&filled_expression)?;
    Ok(regex.is_match(text))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tmpl;
    use std::collections::HashMap;

    #[test]
    fn test_match_expression_plain_text() -> anyhow::Result<()> {
        let text = "Test commit message";
        let expression = tmpl!("^Test");

        assert!(match_expression(&expression, text)?);

        Ok(())
    }

    #[test]
    fn test_match_expression_negative() -> anyhow::Result<()> {
        let text = "Another commit message";
        let expression = tmpl!("^Test.*");

        assert!(!match_expression(&expression, text)?);

        Ok(())
    }

    #[test]
    fn test_match_expression_with_placeholders() -> anyhow::Result<()> {
        let text = "Test commit message";
        let expression = tmpl!(
            "^{{prefix}}",
            HashMap::from([("prefix".to_string(), "Test".to_string())])
        );

        assert!(match_expression(&expression, text)?);

        Ok(())
    }
}
