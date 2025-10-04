'use client'

import { useEffect, useState } from 'react'
import Dashboard from '@/components/Dashboard'
import { FlaggedTransaction, Statistics, WebSocketMessage } from '@/types'

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws'
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'

export default function Home() {
  const [flaggedTransactions, setFlaggedTransactions] = useState<FlaggedTransaction[]>([])
  const [statistics, setStatistics] = useState<Statistics | null>(null)
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    // Fetch initial data
    fetchFlaggedTransactions()
    fetchStatistics()

    // Setup WebSocket connection
    const ws = new WebSocket(WS_URL)

    ws.onopen = () => {
      console.log('WebSocket connected')
      setIsConnected(true)
    }

    ws.onclose = () => {
      console.log('WebSocket disconnected')
      setIsConnected(false)
    }

    ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        
        switch (message.type) {
          case 'fraud_alert':
            // Add new flagged transaction to the top
            fetchFlaggedTransactions()
            break
          
          case 'stats_update':
            setStatistics(message.payload as Statistics)
            break
          
          default:
            console.log('Unknown message type:', message.type)
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }

    return () => {
      ws.close()
    }
  }, [])

  const fetchFlaggedTransactions = async () => {
    try {
      const response = await fetch(`${API_URL}/transactions?limit=50`)
      const data = await response.json()
      setFlaggedTransactions(data || [])
    } catch (error) {
      console.error('Error fetching flagged transactions:', error)
    }
  }

  const fetchStatistics = async () => {
    try {
      const response = await fetch(`${API_URL}/stats`)
      const data = await response.json()
      setStatistics(data)
    } catch (error) {
      console.error('Error fetching statistics:', error)
    }
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <Dashboard
        flaggedTransactions={flaggedTransactions}
        statistics={statistics}
        isConnected={isConnected}
      />
    </main>
  )
}
