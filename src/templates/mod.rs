mod errors;
mod hashmap;
mod vec;
mod template_legacy;
mod tmpl_macro_legacy;
mod template;

pub use crate::templates::errors::TemplateError;
pub use crate::templates::hashmap::replace_in_hashmap;
pub use crate::templates::template_legacy::TemplateStringLegacy;
pub use crate::templates::vec::replace_in_vac;
