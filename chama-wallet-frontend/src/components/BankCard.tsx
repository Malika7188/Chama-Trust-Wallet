import React from 'react';

interface BankCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode;
  progress?: number;
  colorClass?: string;
}

const BankCard: React.FC<BankCardProps> = ({ title, value, icon, progress, colorClass }) => {
  