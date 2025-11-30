import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import App from './App'

// Mock environment variable
vi.stubEnv('VITE_BACKEND_HOST', 'http://localhost:8080')

describe('App', () => {
  it('renders main heading and button', () => {
    render(<App />)
    
    expect(screen.getByRole('main')).toBeInTheDocument()
    expect(screen.getByText("Here's some unnecessary quotes for you to read...")).toBeInTheDocument()
    expect(screen.getByText('Start Quotes')).toBeInTheDocument()
  })

  it('toggles button text when clicked', () => {
    render(<App />)
    
    const button = screen.getByText('Start Quotes')
    fireEvent.click(button)
    
    expect(screen.getByText('Stop Quotes')).toBeInTheDocument()
  })

  it('initially shows no messages', () => {
    render(<App />)
    
    const messages = screen.queryAllByRole('paragraph')
    expect(messages.length).toBe(0)
  })
})