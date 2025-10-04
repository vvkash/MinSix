import { FlaggedTransaction, Statistics } from '@/types'
import { Shield, Activity, AlertTriangle, CheckCircle } from 'lucide-react'
import TransactionList from './TransactionList'
import StatsCard from './StatsCard'
import Header from './Header'

interface DashboardProps {
  flaggedTransactions: FlaggedTransaction[]
  statistics: Statistics | null
  isConnected: boolean
}

export default function Dashboard({ flaggedTransactions, statistics, isConnected }: DashboardProps) {
  const detectionRate = statistics 
    ? ((statistics.total_flagged / Math.max(statistics.total_transactions, 1)) * 100).toFixed(1)
    : '0.0'

  const accuracy = statistics && statistics.total_flagged > 0
    ? (((statistics.total_flagged - statistics.false_positives) / statistics.total_flagged) * 100).toFixed(1)
    : '0.0'

  return (
    <div className="min-h-screen p-8">
      <Header isConnected={isConnected} />
      
      {/* Statistics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatsCard
          title="Total Transactions"
          value={statistics?.total_transactions.toLocaleString() || '0'}
          icon={<Activity className="w-6 h-6" />}
          color="blue"
        />
        <StatsCard
          title="Flagged Transactions"
          value={statistics?.total_flagged.toLocaleString() || '0'}
          icon={<AlertTriangle className="w-6 h-6" />}
          color="red"
        />
        <StatsCard
          title="Detection Rate"
          value={`${detectionRate}%`}
          icon={<Shield className="w-6 h-6" />}
          color="purple"
        />
        <StatsCard
          title="Accuracy"
          value={`${accuracy}%`}
          icon={<CheckCircle className="w-6 h-6" />}
          color="green"
        />
      </div>

      {/* Flagged Transactions List */}
      <div className="bg-white/10 backdrop-blur-lg rounded-xl p-6 shadow-2xl border border-white/20">
        <h2 className="text-2xl font-bold text-white mb-6 flex items-center gap-2">
          <AlertTriangle className="w-6 h-6 text-red-400" />
          Flagged Transactions
        </h2>
        <TransactionList transactions={flaggedTransactions} />
      </div>
    </div>
  )
}
