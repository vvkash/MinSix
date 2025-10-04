export interface Transaction {
  id: number
  tx_hash: string
  block_number: number
  from_address: string
  to_address: string | null
  value: string
  gas_price: string
  gas_used: number
  timestamp: string
}

export interface FlaggedTransaction {
  id: number
  transaction_id: number | null
  tx_hash: string
  risk_score: number
  reasons: string[]
  flagged_at: string
  status: 'pending' | 'reviewed' | 'false_positive' | 'confirmed'
  transaction?: Transaction
}

export interface Statistics {
  total_transactions: number
  total_flagged: number
  false_positives: number
  confirmed_fraud: number
}

export interface WebSocketMessage {
  type: string
  payload: any
}

export interface AlertPayload {
  tx_hash: string
  risk_score: number
  reasons: string[]
  timestamp: string
}
