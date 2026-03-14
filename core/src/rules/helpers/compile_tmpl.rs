use crate::context::Context;
use crate::templates::TemplateString;
use anyhow::Result;

pub fn compile_tmpl(
    ctx: &dyn Context,
    string: &TemplateString,
    additional: &[String],
) -> Result<String> {
    let variables = ctx.variables(additional)?;
    Ok(string.to_string(&variables)?)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_compile_tmpl_success() {
        let mut ctx = MockContext::new();
        let mut vars = HashMap::new();
        vars.insert("name".to_string(), "test".to_string());
        ctx.expect_variables().returning(move |_| Ok(vars.clone()));

        let template = t!("Hello {{name}}");
        let result = compile_tmpl(&ctx, &template, &[]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), "Hello test");
    }

    #[test]
    fn test_compile_tmpl_with_additional() {
        let mut ctx = MockContext::new();
        let mut vars = HashMap::new();
        vars.insert("name".to_string(), "test".to_string());
        ctx.expect_variables().returning(move |_| Ok(vars.clone()));

        let template = t!("Hello {{name}}");
        let result = compile_tmpl(&ctx, &template, &["extra".to_string()]);
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), "Hello test");
    }

    #[test]
    fn test_compile_tmpl_variables_error() {
        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let template = t!("Hello {{name}}");
        let result = compile_tmpl(&ctx, &template, &[]);
        assert!(result.is_err());
    }

    #[test]
    fn test_compile_tmpl_template_error() {
        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let template = t!("Hello {{missing}}");
        let result = compile_tmpl(&ctx, &template, &[]);
        assert!(result.is_err());
    }
}
