mod placeholders;
mod errors;
mod hashmap;
mod vec;
mod template;
mod tmpl_macro;

pub use crate::templates::errors::TemplateError;
pub use crate::templates::hashmap::replace_in_hashmap;
pub use crate::templates::placeholders::replace_in_string;
pub use crate::templates::vec::replace_in_vac;
pub use crate::templates::template::TemplateString;
