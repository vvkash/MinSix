import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Minsix - Ethereum Fraud Detection',
  description: 'Real-time fraud detection platform for Ethereum wallets',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
