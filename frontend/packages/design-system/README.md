# Volha design system

Native React components styled with SCSS Modules. Components use the bundled Volha theme and do not depend on Radix UI, CVA, or Tailwind utility classes.

Import the bundled theme once in the application entry point or root layout:

```tsx
import '../src/shared/ui/design-system/styles/theme.scss'
```

The theme exposes portable `--volha-*` CSS custom properties. Sass consumers can also access the source maps and mixins through `styles/index.scss`.

```tsx
import {
  BackLink,
  Button,
  InputField,
  InstructionSteps,
  Modal,
  ProjectMeta,
  Switch,
  ThumbnailGallery,
} from '../src/shared/ui/design-system'
```

## Groups

- Actions: `Button`, `IconButton`, `Chip`, `SegmentedControl`
- Navigation: `BackLink`
- Forms: `InputField`, `TextareaField`, `SelectField`, `Checkbox`, `RadioGroup`, `Switch`
- Feedback: `Badge`, `Alert`, `Spinner`
- Overlay: `Modal`
- Commerce: `ProductCard`, `VolhaLogo`
- Project patterns: `ProjectMeta`, `InstructionSteps`, `ThumbnailGallery`

`Modal` includes focus trapping, Escape handling, scroll locking, backdrop closing, focus restoration, and responsive bottom-sheet behavior.

Components use the project dependency `clsx` for conditional class names.

## DIY project patterns

`Badge` includes the semantic tones `easy`, `medium`, `hard`, and `hit`. The
matching colors are bundled into the theme as portable `--volha-color-*`
variables.

```tsx
<BackLink href="/diy">Все проекты</BackLink>

<ProjectMeta
  category="Гараж"
  difficulty="Средний"
  duration="2 дня"
/>

<InstructionSteps
  steps={[
    { label: 'Замеры и планировка' },
    { label: 'Монтаж стеллажей S-200' },
  ]}
/>
```

Use `ProductCard variant="project"` for the quieter product tiles used inside
DIY articles. This variant removes the commerce divider and keeps price
metadata compact.
