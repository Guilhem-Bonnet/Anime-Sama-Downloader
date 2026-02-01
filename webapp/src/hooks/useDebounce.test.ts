import { describe, it, expect, vi } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { useDebounce } from '../useDebounce'

describe('useDebounce Hook', () => {
  it('should return initial value immediately', () => {
    const { result } = renderHook(() => useDebounce('test', 300))
    expect(result.current).toBe('test')
  })

  it('should debounce value changes', async () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 100),
      { initialProps: { value: 'initial' } }
    )

    expect(result.current).toBe('initial')

    rerender({ value: 'updated' })

    // Value shouldn't update immediately
    expect(result.current).toBe('initial')

    // Wait for debounce delay
    await waitFor(
      () => {
        expect(result.current).toBe('updated')
      },
      { timeout: 200 }
    )
  })

  it('should handle custom delay', async () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 50),
      { initialProps: { value: 'test1' } }
    )

    rerender({ value: 'test2' })
    expect(result.current).toBe('test1')

    await waitFor(
      () => {
        expect(result.current).toBe('test2')
      },
      { timeout: 150 }
    )
  })

  it('should handle multiple rapid changes', async () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 100),
      { initialProps: { value: 'a' } }
    )

    rerender({ value: 'ab' })
    rerender({ value: 'abc' })
    rerender({ value: 'abcd' })

    await waitFor(
      () => {
        expect(result.current).toBe('abcd')
      },
      { timeout: 200 }
    )
  })

  it('should work with different data types', async () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 50),
      { initialProps: { value: { count: 0 } } }
    )

    rerender({ value: { count: 1 } })

    await waitFor(
      () => {
        expect(result.current.count).toBe(1)
      },
      { timeout: 150 }
    )
  })

  it('should work with arrays', async () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 50),
      { initialProps: { value: [] } }
    )

    const newArray = [1, 2, 3]
    rerender({ value: newArray })

    await waitFor(
      () => {
        expect(result.current).toEqual(newArray)
      },
      { timeout: 150 }
    )
  })
})
