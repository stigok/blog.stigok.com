import unittest
import pymd

SCRIPT_NO_EVAL = """
start
```python
print(1)
```
end
"""
SCRIPT_EVAL = """
start
```python
print(1)
#eval
```
end
"""


class TestPymd(unittest.TestCase):
    def test_no_eval(self):
        res = pymd.process_text(SCRIPT_NO_EVAL)
        self.assertEqual(
            res,
            """
start
```python
print(1)
```
end
""",
        )

    def test_simple_script(self):
        res = pymd.process_text(SCRIPT_EVAL)
        self.assertEqual(
            res,
            """
start
```python
print(1)
# 1
```
end
""",
        )

    def test_one_normal_one_eval(self):
        """
        Make sure it's okay to have a non-eval script before an eval script
        within the same text, and that only the eval one gets evaled.
        """
        text = SCRIPT_NO_EVAL + SCRIPT_EVAL
        res = pymd.process_text(text)
        self.assertEqual(
            res,
            """
start
```python
print(1)
```
end

start
```python
print(1)
# 1
```
end
""",
        )
