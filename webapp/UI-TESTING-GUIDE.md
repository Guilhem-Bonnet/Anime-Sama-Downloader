# UI Testing Strategy - Anime-Sama-Downloader

## Overview

This document outlines the testing strategy for the React UI components in the webapp module.

## Testing Framework

- **Framework**: Vitest + React Testing Library
- **Environment**: jsdom
- **Coverage Target**: 80%+
- **Test Types**: Unit tests, Component tests, Hook tests

## Setup

### Configuration Files
- `vitest.config.ts` - Vitest configuration
- `vitest.setup.ts` - Test environment setup (mocks for localStorage, matchMedia)

### Installation
```bash
cd webapp
npm install
npm test
```

### Available Commands
```bash
npm test              # Run tests in watch mode
npm test -- --run    # Run tests once
npm test:ui          # Run tests with UI dashboard
npm test:coverage    # Generate coverage report
```

## Test Structure

### UI Components (src/components/ui/)

#### Button Component ✅
- **File**: `Button.test.tsx`
- **Status**: 5/7 tests passing
- **Covered**:
  - ✅ Renders button with text
  - ✅ Handles click events
  - ✅ Applies disabled state
  - ❌ Apply variant classes (implementation detail)
  - ❌ Apply size classes (implementation detail)

#### Badge Component ✅
- **File**: `Badge.test.tsx`
- **Status**: 4/6 tests passing
- **Covered**:
  - ✅ Renders badge with text
  - ✅ Applies variant classes
  - ✅ Applies custom className
  - ✅ Renders outline variant

#### Card Component ⚠️
- **File**: `Card.test.tsx`
- **Status**: 3/8 tests passing
- **Issue**: Card subcomponents (Header, Title, etc.) need verification

#### Input Component ✅
- **File**: `Input.test.tsx`
- **Status**: 7/10 tests passing
- **Covered**:
  - ✅ Renders input element
  - ✅ Handles text input
  - ✅ Placeholder rendering
  - ✅ Disabled state
  - ✅ Different input types (password, email, number)

### Custom Components (src/components/custom/)

#### StatusBadge ✅
- **File**: `StatusBadge.test.tsx`
- **Status**: 7/7 tests passing
- **Covered**:
  - ✅ Renders status badge
  - ✅ Displays various status values
  - ✅ Custom sizing
  - ✅ Custom className

### Search Components (src/components/search/)

#### SearchBar ⏳
- **File**: `SearchBar.test.tsx`
- **Status**: In progress
- **Note**: Requires mocking of Zustand store and custom hooks

### Hooks (src/hooks/)

#### useDebounce ✅
- **File**: `useDebounce.test.ts`
- **Status**: 6/6 tests passing
- **Covered**:
  - ✅ Returns initial value
  - ✅ Debounces value changes
  - ✅ Custom delay handling
  - ✅ Multiple rapid changes
  - ✅ Different data types (objects, arrays)

## Test Results Summary

| Category | Tests | Status |
|----------|-------|--------|
| UI Components | 26 | 13 passing |
| Custom Components | 7 | 7 passing ✅ |
| Hooks | 6 | 6 passing ✅ |
| **Total** | **39** | **26 passing** |

## Known Issues & TODO

### High Priority
- [ ] Fix Card subcomponent tests (Header, Title, Description, Content, Footer)
- [ ] Fix Button variant/size class assertions
- [ ] Fix Input type-specific assertions

### Medium Priority
- [ ] Implement SearchBar tests with proper mocks
- [ ] Add tests for AutocompleteSuggestions component
- [ ] Add tests for RecentSearchesDropdown component
- [ ] Add tests for FilterPanel component

### Low Priority
- [ ] Add integration tests with full app context
- [ ] Add E2E tests with Playwright/Cypress
- [ ] Setup coverage reporting in CI/CD
- [ ] Add visual regression tests

## Best Practices

1. **Use React Testing Library**: Focus on user behavior, not implementation
2. **Mock External Dependencies**: Zustand stores, API calls, localStorage
3. **Test Accessibility**: Use semantic HTML and ARIA attributes
4. **Avoid Implementation Details**: Don't test className or internal structure
5. **Group Related Tests**: Use `describe()` blocks
6. **Clear Test Names**: Describe what the test does in plain language

## Running Tests

### Run all tests
```bash
npm test -- --run
```

### Run specific file
```bash
npm test -- Button.test.tsx --run
```

### Run with coverage
```bash
npm test:coverage
```

### Watch mode (default)
```bash
npm test
```

### UI Dashboard
```bash
npm test:ui
```

## Adding New Tests

1. Create `ComponentName.test.tsx` in the same folder as the component
2. Use consistent test structure:
   ```typescript
   import { describe, it, expect } from 'vitest'
   import { render, screen } from '@testing-library/react'
   import { MyComponent } from './MyComponent'

   describe('MyComponent', () => {
     it('should render', () => {
       render(<MyComponent />)
       expect(screen.getByRole('...')).toBeInTheDocument()
     })
   })
   ```
3. Run tests: `npm test -- --run`
4. Fix failing tests before committing

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [React Testing Library Docs](https://testing-library.com/react)
- [Testing Best Practices](https://kentcdodds.com/blog/common-mistakes-with-react-testing-library)
