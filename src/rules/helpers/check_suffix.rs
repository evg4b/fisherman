use crate::context::Context;
use crate::rules::helpers::compile_tmpl::compile_tmpl;
use crate::templates::TemplateString;

pub fn check_suffix(
    ctx: &dyn Context,
    suffix: &TemplateString,
    text: &str,
) -> anyhow::Result<bool> {
    Ok(text.ends_with(&compile_tmpl(ctx, suffix, &[])?))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_check_suffix() -> anyhow::Result<()> {
        let text = "Test commit message";
        let suffix = t!("message");
        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        assert!(check_suffix(&context, &suffix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_suffix_negative() -> anyhow::Result<()> {
        let text = "Test commit message";
        let suffix = t!("commit");
        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        assert!(!check_suffix(&context, &suffix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_suffix_with_placeholders() -> anyhow::Result<()> {
        let text = "Test commit message";
        let suffix = t!("{{suffix}}");
        let mut context = MockContext::new();
        context.expect_variables().returning(|_| {
            Ok(HashMap::from([(
                "suffix".to_string(),
                "message".to_string(),
            )]))
        });

        assert!(check_suffix(&context, &suffix, text)?);

        Ok(())
    }
}
