extern crate proc_macro;

use proc_macro::TokenStream;
use quote::quote;
use syn;
use syn::{parse_macro_input, DeriveInput};

#[proc_macro_derive(ConditionalRule)]
pub fn conditional_rule_derive(input: TokenStream) -> TokenStream {
    let input = parse_macro_input!(input as DeriveInput);
    let name = &input.ident;

    let sss = quote! {
        impl ConditionalRule for #name {
            fn check_condition(&self, ctx: &dyn Context) -> anyhow::Result<bool> {
                self.when.as_ref().map(|expr| expr.check(ctx.variables(&[]))).unwrap_or(Ok(false))
            }
        }
    };

    sss.into()
}