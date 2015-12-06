pub struct Foo {
  value:uint
}

trait HasUIntValue {
  fn as_uint(self) -> uint;
}

pub impl Foo {
  fn add<T:HasUIntValue>(&mut self, value:T) {
    self.value += value.as_uint();
  }
}

impl HasUIntValue for int {
  fn as_uint(self) -> uint {
    return self as uint;
  }
}

impl HasUIntValue for (f64) {
  fn as_uint(self) -> uint {
    return self as uint;
  }
}

#[test]
fn test_add_with_int()
{
  let mut x = Foo { value: 10 };
  x.add(10i);
  assert!(x.value == 21);
}

#[test]
fn test_add_with_float()
{
  let mut x = Foo { value: 10 };
  x.add(10.0f64);
  assert!(x.value == 20);
}

