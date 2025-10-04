import { FlaggedTransaction } from '@/types'
import { ExternalLink, AlertCircle } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'

interface TransactionCardProps {
  transaction: FlaggedTransaction
}

export default function TransactionCard({ transaction }: TransactionCardProps) {
  const getRiskColor = (score: number) => {
    if (score >= 70) return 'text-red-400 bg-red-500/20'
    if (score >= 40) return 'text-orange-400 bg-orange-500/20'
    return 'text-yellow-400 bg-yellow-500/20'
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'confirmed':
        return 'bg-red-500/20 text-red-400'
      case 'false_positive':
        return 'bg-green-500/20 text-green-400'
      case 'reviewed':
        return 'bg-blue-500/20 text-blue-400'
      default:
        return 'bg-gray-500/20 text-gray-400'
    }
  }

  const etherscanUrl = `https://etherscan.io/tx/${transaction.tx_hash}`
  const timeAgo = formatDistanceToNow(new Date(transaction.flagged_at), { addSuffix: true })

  const formatValue = (value: string) => {
    try {
      const eth = parseFloat(value) / 1e18
      return eth.toFixed(4)
    } catch {
      return '0.0000'
    }
  }

  return (
    <div className="bg-white/5 backdrop-blur-sm rounded-lg p-4 border border-white/10 hover:bg-white/10 transition-all">
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <a
              href={etherscanUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="text-purple-400 hover:text-purple-300 font-mono text-sm flex items-center gap-1"
            >
              {transaction.tx_hash.slice(0, 10)}...{transaction.tx_hash.slice(-8)}
              <ExternalLink className="w-4 h-4" />
            </a>
            <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(transaction.status)}`}>
              {transaction.status}
            </span>
          </div>
          
          {transaction.transaction && (
            <div className="text-sm text-gray-400 space-y-1">
              <div className="flex items-center gap-2">
                <span className="text-gray-500">From:</span>
                <span className="font-mono">{transaction.transaction.from_address.slice(0, 8)}...{transaction.transaction.from_address.slice(-6)}</span>
              </div>
              {transaction.transaction.to_address && (
                <div className="flex items-center gap-2">
                  <span className="text-gray-500">To:</span>
                  <span className="font-mono">{transaction.transaction.to_address.slice(0, 8)}...{transaction.transaction.to_address.slice(-6)}</span>
                </div>
              )}
              <div className="flex items-center gap-2">
                <span className="text-gray-500">Value:</span>
                <span className="font-medium text-white">{formatValue(transaction.transaction.value)} ETH</span>
              </div>
            </div>
          )}
        </div>

        <div className={`px-3 py-2 rounded-lg font-bold text-lg ${getRiskColor(transaction.risk_score)}`}>
          {transaction.risk_score}
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-start gap-2">
          <AlertCircle className="w-4 h-4 text-red-400 mt-0.5 flex-shrink-0" />
          <div className="flex-1">
            <p className="text-xs text-gray-500 mb-1">Reasons:</p>
            <div className="flex flex-wrap gap-2">
              {transaction.reasons.map((reason, index) => (
                <span
                  key={index}
                  className="px-2 py-1 bg-red-500/20 text-red-300 rounded text-xs"
                >
                  {reason}
                </span>
              ))}
            </div>
          </div>
        </div>

        <div className="text-xs text-gray-500 text-right">
          {timeAgo}
        </div>
      </div>
    </div>
  )
}
