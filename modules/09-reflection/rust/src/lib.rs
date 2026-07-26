//! Module 09 in Rust: the answer to "operate generically over a type" is
//! COMPILE-TIME CODE GENERATION via `#[derive]`. A derive macro is a procedural
//! macro that runs at compile time and emits a concrete `impl` block for exactly
//! this type -- no runtime type descriptor, no reflection walk. The std derives
//! (Debug, Clone, PartialEq, Hash, Default) are built in; serde's Serialize /
//! Deserialize are the same mechanism supplied by a crate. Errors are caught at
//! compile time: a derive fails to compile if a field does not satisfy the
//! bound (see reject-rust), the opposite of reflection's deferral to run time.

// region:derive:start

// The compiler generates Debug (a formatter), Clone, and PartialEq (a
// field-by-field ==) for this exact type. Each is real generated code, checked
// against every field's type at compile time.
#[derive(Debug, Clone, PartialEq)]
pub struct Person {
    pub name: String,
    pub age: u32,
}

// region:derive:end

impl Person {
    pub fn new(name: &str, age: u32) -> Person {
        Person {
            name: name.to_string(),
            age,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn derived_impls() {
        let a = Person::new("Ada", 36);
        let b = a.clone(); // derived Clone
        assert_eq!(a, b); // derived PartialEq
        assert_eq!(format!("{a:?}"), r#"Person { name: "Ada", age: 36 }"#); // derived Debug
        assert_ne!(a, Person::new("Ada", 37));
    }
}
