import { FlaggedTransaction } from '@/types'
import TransactionCard from './TransactionCard'

interface TransactionListProps {
  transactions: FlaggedTransaction[]
}

export default function TransactionList({ transactions }: TransactionListProps) {
  if (transactions.length === 0) {
    return (
      <div className="text-center py-12 text-gray-400">
        <p className="text-lg">No flagged transactions yet.</p>
        <p className="text-sm mt-2">The system is monitoring for suspicious activity.</p>
      </div>
    )
  }

  return (
    <div className="space-y-4 max-h-[600px] overflow-y-auto pr-2">
      {transactions.map((transaction) => (
        <TransactionCard key={transaction.id} transaction={transaction} />
      ))}
    </div>
  )
}
