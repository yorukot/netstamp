import { useState } from 'react'
import { type AuthCredentials, type TeamDraft, getSessionSnapshot, mockCreateTeam, mockLogin, mockRegister } from '../services/authService'

export function useAuthMock() {
  const [session, setSession] = useState(() => getSessionSnapshot())
  const [submitting, setSubmitting] = useState(false)

  async function login(payload: AuthCredentials) {
    setSubmitting(true)
    const user = await mockLogin(payload)
    setSession({ user, controller: 'waiting-for-api' })
    setSubmitting(false)
    return user
  }

  async function register(payload: AuthCredentials) {
    setSubmitting(true)
    const user = await mockRegister(payload)
    setSession({ user, controller: 'waiting-for-api' })
    setSubmitting(false)
    return user
  }

  async function createTeam(payload: TeamDraft) {
    setSubmitting(true)
    const team = await mockCreateTeam(payload)
    setSession((current) => ({ ...current, team }))
    setSubmitting(false)
    return team
  }

  return { session, submitting, login, register, createTeam }
}
