# Safe Strings Composite Action

A helper action to convert a string to only contain alphanumeric values and the option to limit its overall length. Intented to be used for tagging and similar activities where other characters (like `/` and `|`) cause errors.

String is converted to lowercase version.

## Usage

Within you github workflow job you can place a step such as this:

```yaml
    - name: "Generate safe string"
      id: safe_string
      uses: 'ministryofjustice/opg-github-actions/.github/actions/branch-name@v2.3.1'
      with:
        original: "none-safe-string"
```

## Inputs and Outputs

Inputs:
- `original`
- `length`
- `suffix`
- `conditional_match`
- `conditional_value`


Outputs:
- `original`
- `length`
- `suffix`
- `full_length`
- **`safe`**


### Inputs

#### `original`
Your original string that you want to convert over be safe for use in tags and similar.

#### `length`
If set, `safe` version of `original` will also be truncated to this length.

#### `suffix`
An option to allow a suffix to be added to the cleansed string. This value is not cleaned beforehand.

If length has been set as well then this will be included in the max length string - eg `original` is set as `test123`, `length` is `6` and `suffix` is `AB` the result would be `testAB`

#### `conditional_match` and `conditional_value`
Configuring these allows an overwrite of the final resulting string. This is to allow certain inputs to force a set result, for example, if `original` equals `main` and you set `conditional_match` to 'main' and `conditional_value` to `production` then `safe` would be set to `production`.

This designed to allow flexibility in naming of production or other fixed environments, such as 'dev' => 'development' etc.

### Outputs

#### `original`
Mirror of the input value.

#### `length`
Mirror of the input value.

#### `suffix`
Mirror of the input value.

#### `full_length`
A santised version of `original`, but without any truncation.

#### `safe`
A santised version of `original`, that has also been truncated to a set length.
