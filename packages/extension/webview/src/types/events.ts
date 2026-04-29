export type EventType =
  | 'log'
  | 'error'
  | 'warning'
  | 'process_start'
  | 'process_exit'
  | 'stage_change'

export type StageType = 'build' | 'run' | 'test' | 'unknown'

export type StageStatus = 'idle' | 'running' | 'success' | 'error'

// ── Payloads por tipo de evento ──────────────────────────────────────────────

export interface LogPayload {
  message: string
  source: 'stdout' | 'stderr'
}

export interface ErrorPayload {
  message: string
  source: 'stdout' | 'stderr'
}

export interface WarningPayload {
  message: string
  source: 'stdout' | 'stderr'
}

export interface ProcessStartPayload {
  bin:  string
  args: string[]
  pid:  number
}

export interface ProcessExitPayload {
  exitCode:   number
  durationMs: number
}

export interface StageChangePayload {
  previousStage: StageType
  newStage:      StageType
}

// ── Tipos discriminados por EventType ────────────────────────────────────────

export interface LogEvent {
  id:        string
  type:      'log'
  stage:     StageType
  timestamp: number
  seq:       number
  payload:   LogPayload
}

export interface ErrorEvent {
  id:        string
  type:      'error'
  stage:     StageType
  timestamp: number
  seq:       number
  payload:   ErrorPayload
}

export interface WarningEvent {
  id:        string
  type:      'warning'
  stage:     StageType
  timestamp: number
  seq:       number
  payload:   WarningPayload
}

export interface ProcessStartEvent {
  id:        string
  type:      'process_start'
  stage:     StageType
  timestamp: number
  seq:       number
  payload:   ProcessStartPayload
}

export interface ProcessExitEvent {
  id:        string
  type:      'process_exit'
  stage:     StageType
  timestamp: number
  seq:       number
  payload:   ProcessExitPayload
}

export interface StageChangeEvent {
  id:        string
  type:      'stage_change'
  stage:     StageType
  timestamp: number
  seq:       number
  payload:   StageChangePayload
}

export type ExrayEvent =
  | LogEvent
  | ErrorEvent
  | WarningEvent
  | ProcessStartEvent
  | ProcessExitEvent
  | StageChangeEvent

// ── Comando enviado da extensão para o Agent Go via stdin ────────────────────

export interface AgentCommand {
  command:  'run' | 'stop'
  bin?:     string
  args?:    string[]
  stage?:   StageType
  timeout?: number
}