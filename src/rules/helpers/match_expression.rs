use crate::context::Context;
use crate::templates::TemplateString;
use regex::Regex;

pub fn match_expression(
    ctx: &dyn Context,
    expression: &TemplateString,
    text: &str,
) -> anyhow::Result<bool> {
    let variables = ctx.variables(&[])?;
    let filled_expression = expression.to_string(&variables)?;
    let regex = Regex::new(&filled_expression)?;
    Ok(regex.is_match(text))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_match_expression_plain_text() -> anyhow::Result<()> {
        let text = "Test commit message";
        let expression = t!("^Test");

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        assert!(match_expression(&context, &expression, text)?);

        Ok(())
    }

    #[test]
    fn test_match_expression_negative() -> anyhow::Result<()> {
        let text = "Another commit message";
        let expression = t!("^Test.*");

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        assert!(!match_expression(&context, &expression, text)?);

        Ok(())
    }

    #[test]
    fn test_match_expression_with_placeholders() -> anyhow::Result<()> {
        let text = "Test commit message";
        let expression = t!("^{{prefix}}");

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::from([("prefix".to_string(), "Test".to_string())])));

        assert!(match_expression(&context, &expression, text)?);

        Ok(())
    }
}
