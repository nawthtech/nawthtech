import '@testing-library/jest-dom'
import { expect, afterEach, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import * as matchers from '@testing-library/jest-dom/matchers'

expect.extend(matchers)

// Mock window.scrollTo
Object.defineProperty(window, 'scrollTo', {
  value: vi.fn(),
  writable: true
})

// Mock EventSource
class MockEventSource {
  static readonly CONNECTING = 0
  static readonly OPEN = 1
  static readonly CLOSED = 2

  public onopen: (() => void) | null = null
  public onmessage: ((event: any) => void) | null = null
  public onerror: ((event: any) => void) | null = null
  public readyState = MockEventSource.CONNECTING
  public url: string

  constructor(url: string) {
    this.url = url
    // Simulate connection opening
    setTimeout(() => {
      this.readyState = MockEventSource.OPEN
      if (this.onopen) this.onopen()
    }, 0)
  }

  close() {
    this.readyState = MockEventSource.CLOSED
  }

  // Helper method for tests to simulate messages
  simulateMessage(data: string) {
    if (this.onmessage && this.readyState === MockEventSource.OPEN) {
      this.onmessage({ data })
    }
  }

  // Helper method for tests to simulate errors
  simulateError() {
    if (this.onerror) {
      this.onerror(new Event('error'))
      this.readyState = MockEventSource.CLOSED
    }
  }
}

// @ts-ignore
global.EventSource = MockEventSource

// Mock environment variables
vi.stubEnv('VITE_BACKEND_HOST', 'http://localhost:8080')

afterEach(() => {
  cleanup()
  vi.clearAllMocks()
})