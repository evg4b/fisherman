use crate::templates::TemplateString;

pub fn check_prefix(prefix: &TemplateString, text: &str) -> anyhow::Result<bool> {
    let filled_prefix = prefix.to_string()?;
    Ok(text.starts_with(&filled_prefix))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tmpl;
    use std::collections::HashMap;

    #[test]
    fn test_check_prefix() -> anyhow::Result<()> {
        let text = "Test commit message";
        let prefix = tmpl!("Test");

        assert!(check_prefix(&prefix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_prefix_negative() -> anyhow::Result<()> {
        let text = "Another commit message";
        let prefix = tmpl!("Test");

        assert!(!check_prefix(&prefix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_prefix_with_placeholders() -> anyhow::Result<()> {
        let text = "Test commit message";
        let prefix = tmpl!(
            "{{prefix}}",
            HashMap::from([("prefix".to_string(), "Test".to_string())])
        );

        assert!(check_prefix(&prefix, text)?);

        Ok(())
    }
}