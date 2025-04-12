mod errors;
mod hashmap;
mod vec;
mod template;
mod tmpl_macro;

pub use crate::templates::errors::TemplateError;
pub use crate::templates::hashmap::replace_in_hashmap;
pub use crate::templates::template::TemplateString;
pub use crate::templates::vec::replace_in_vac;
