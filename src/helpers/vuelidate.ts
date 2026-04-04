interface ErrorObject {
  $message: string
}

interface ValidationField {
  $error: boolean
  $errors: ErrorObject[]
}

export const getErrorMessage = (field: ValidationField): any => {
  if (!field.$error || field.$errors.length === 0) return ''
  return field.$errors[0].$message
}
