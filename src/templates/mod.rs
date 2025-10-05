mod errors;
mod hashmap;
mod template;
mod template_macro;
mod vec;

pub use crate::templates::errors::TemplateError;
pub use crate::templates::hashmap::replace_in_hashmap;
pub use crate::templates::template::TemplateString;
pub use crate::templates::vec::replace_in_vec;
