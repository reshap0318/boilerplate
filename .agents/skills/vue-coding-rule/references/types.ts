// ============================================================
// Types File — grouped by component
// ============================================================

// UiCard types
export interface UiCardClasses {
  wrapper?: string
  card?: string
  header?: string
  body?: string
  footer?: string
}

export interface UiCardProps {
  title?: string
  classes?: UiCardClasses
}

// UiModal types
export interface UiModalProps {
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | '2xl'
}

// UiButton types
export interface UiButtonProps {
  label?: string
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  disabled?: boolean
  loading?: boolean
}
