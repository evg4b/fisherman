#[macro_export] macro_rules! tmpl {
    // Case 1: string literal only, empty HashMap
    ($template:expr) => {
        $crate::templates::TemplateString::new($template.to_string(), std::collections::HashMap::new())
    };

    // Case 2: string literal with variables
    ($template:expr, $variables:expr) => {
        $crate::templates::TemplateString::new($template.to_string(), $variables)
    };
}