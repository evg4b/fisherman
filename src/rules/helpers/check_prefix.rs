use crate::templates::TemplateStringLegacy;

pub fn check_prefix(prefix: &TemplateStringLegacy, text: &str) -> anyhow::Result<bool> {
    let filled_prefix = prefix.to_string()?;
    Ok(text.starts_with(&filled_prefix))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tmpl_legacy;
    use std::collections::HashMap;

    #[test]
    fn test_check_prefix() -> anyhow::Result<()> {
        let text = "Test commit message";
        let prefix = tmpl_legacy!("Test");

        assert!(check_prefix(&prefix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_prefix_negative() -> anyhow::Result<()> {
        let text = "Another commit message";
        let prefix = tmpl_legacy!("Test");

        assert!(!check_prefix(&prefix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_prefix_with_placeholders() -> anyhow::Result<()> {
        let text = "Test commit message";
        let prefix = tmpl_legacy!(
            "{{prefix}}",
            HashMap::from([("prefix".to_string(), "Test".to_string())])
        );

        assert!(check_prefix(&prefix, text)?);

        Ok(())
    }
}