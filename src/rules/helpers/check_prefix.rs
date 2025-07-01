use crate::context::Context;
use crate::rules::helpers::compile_tmpl::compile_tmpl;
use crate::templates::TemplateString;

pub fn check_prefix(
    ctx: &dyn Context,
    prefix: &TemplateString,
    text: &str,
) -> anyhow::Result<bool> {
    Ok(text.starts_with(&compile_tmpl(ctx, prefix, &[])?))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use std::collections::HashMap;

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
