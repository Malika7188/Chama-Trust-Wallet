import React from 'react';

interface BankCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode;
  progress?: number;
  colorClass?: string;
}

const BankCard: React.FC<BankCardProps> = ({ title, value, icon, progress, colorClass }) => {
  return (
    <div
      className="relative rounded-3xl p-10 flex flex-col gap-4 min-h-[180px] min-w-[340px] max-w-[420px] w-full bg-white/80 backdrop-blur-lg shadow-2xl border border-gray-100 transition-transform duration-300 ease-out hover:-translate-y-3 hover:shadow-3xl group overflow-hidden"
      style={{
        boxShadow: '0 12px 40px 0 rgba(31, 38, 135, 0.18)',
      