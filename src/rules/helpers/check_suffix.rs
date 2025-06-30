use crate::templates::TemplateStringLegacy;

pub fn check_suffix(suffix: &TemplateStringLegacy, text: &str) -> anyhow::Result<bool> {
    let filled_suffix = suffix.to_string()?;
    Ok(text.ends_with(&filled_suffix))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::tmpl_legacy;
    use std::collections::HashMap;

    #[test]
    fn test_check_suffix() -> anyhow::Result<()> {
        let text = "Test commit message";
        let suffix = tmpl_legacy!("message");

        assert!(check_suffix(&suffix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_suffix_negative() -> anyhow::Result<()> {
        let text = "Test commit message";
        let suffix = tmpl_legacy!("commit");

        assert!(!check_suffix(&suffix, text)?);

        Ok(())
    }

    #[test]
    fn test_check_suffix_with_placeholders() -> anyhow::Result<()> {
        let text = "Test commit message";
        let suffix = tmpl_legacy!(
            "{{suffix}}",
            HashMap::from([("suffix".to_string(), "message".to_string())])
        );

        assert!(check_suffix(&suffix, text)?);

        Ok(())
    }
}
