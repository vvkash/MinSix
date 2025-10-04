import { ReactNode } from 'react'

interface StatsCardProps {
  title: string
  value: string
  icon: ReactNode
  color: 'blue' | 'red' | 'green' | 'purple'
}

const colorClasses = {
  blue: 'from-blue-500 to-cyan-500',
  red: 'from-red-500 to-pink-500',
  green: 'from-green-500 to-emerald-500',
  purple: 'from-purple-500 to-indigo-500',
}

export default function StatsCard({ title, value, icon, color }: StatsCardProps) {
  return (
    <div className="bg-white/10 backdrop-blur-lg rounded-xl p-6 shadow-xl border border-white/20 hover:bg-white/15 transition-all">
      <div className="flex items-center justify-between mb-4">
        <div className={`bg-gradient-to-br ${colorClasses[color]} p-3 rounded-lg`}>
          {icon}
        </div>
      </div>
      <h3 className="text-gray-300 text-sm font-medium mb-1">{title}</h3>
      <p className="text-3xl font-bold text-white">{value}</p>
    </div>
  )
}
