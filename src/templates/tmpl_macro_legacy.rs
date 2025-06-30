#[macro_export] macro_rules! tmpl_legacy {
    // Case 1: string literal only, empty HashMap
    ($template:expr) => {
        $crate::templates::TemplateStringLegacy::new($template.to_string(), std::collections::HashMap::new())
    };

    // Case 2: string literal with variables
    ($template:expr, $variables:expr) => {
        $crate::templates::TemplateStringLegacy::new($template.to_string(), $variables)
    };
}