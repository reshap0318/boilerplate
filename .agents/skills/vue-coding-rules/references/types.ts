// ============================================================
// Types File — grouped by component
// Real source of truth: src/components/utils/types.ts — this file
// only illustrates the PATTERN (one interface set per component,
// `classes` object prop for styling hooks). It is NOT the full
// catalog. Always read the real types.ts + list src/components/utils/
// before adding a new component or assuming a prop shape.
// ============================================================

// UiButton types
export interface UiButtonProps {
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  variant?: 'primary' | 'secondary' | 'danger' | 'success'
  outline?: boolean
  size?: 'sm' | 'md' | 'lg'
  rounded?: 'none' | 'sm' | 'md' | 'lg' | 'full'
  fullWidth?: boolean
  loadingText?: string
}

// UiCard types — shows the `classes` object-prop pattern (see SKILL.md
// "Classes Props Pattern"): styling hooks are grouped under one prop
// instead of many individual `xClass` props.
export interface UiCardClasses {
  wrapper?: string
  card?: string
  header?: string
  body?: string
  footer?: string
}

export interface UiCardProps {
  padded?: boolean
  classes?: UiCardClasses
}

// ...remaining components (UiModal, UiPagination, UiSkeleton, UiTable,
// UiEmptyState, FormInput, FormSelect, FormPassword, FormAvatar,
// FormFile, UiBadge, etc.) follow the same one-interface-set-per-
// component shape — check the real file for their current props.
