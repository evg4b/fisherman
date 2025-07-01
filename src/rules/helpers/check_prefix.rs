use crate::context::Context;
use crate::templates::TemplateString;

pub fn check_prefix(
    ctx: &dyn Context,
    prefix: &TemplateString,
    text: &str,
) -> anyhow::Result<bool> {
    let variables = ctx.variables(&[])?;
    let filled_prefix = prefix.to_string(&variables)?;
    Ok(text.starts_with(&filled_prefix))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use std::collections::HashMap;
    use crate::t;

    #[test]
    fn test_check_prefix() -> anyhow::Result<()> {
        let text = "Test commit message";
        let prefix = t!("Test");


        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        assert!(check_prefix(&context, &prefix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_prefix_negative() -> anyhow::Result<()> {
        let text = "Another commit message";
        let prefix = t!("Test");

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        assert!(!check_prefix(&context, &prefix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_prefix_with_placeholders() -> anyhow::Result<()> {
        let text = "Test commit message";
        let prefix = t!("{{prefix}}");

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::from([("prefix".to_string(), "Test".to_string())])));

        assert!(check_prefix(&context, &prefix, text)?);

        Ok(())
    }
}
