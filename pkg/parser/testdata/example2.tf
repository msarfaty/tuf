resource "aws_security_group" "foobar" {
  count = 1

  name = "some_sg_group"
}
