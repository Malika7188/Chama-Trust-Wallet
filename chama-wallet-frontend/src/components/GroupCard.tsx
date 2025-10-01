import { Link } from 'react-router-dom';
import { Users } from 'lucide-react';
import { useGroupBalance } from '../hooks/useGroups';
import type { Group } from '../types';

interface GroupCardProps {
  group: Group;
}

const GroupCard = ({ group }: GroupCardProps) => {
  const { data: balance, isLoading: balanceLoading } = useGroupBalance(group.ID);

  const groupBalance = balance?.data?.balance || '0';

  return (
    <Link
      key={group.ID}
      to={`/groups/${group.ID}`}
      className="card hover:shadow-md transition-shadow"
    >
      <div className="flex items-center justify-between mb-4">
        <div className="w-12 h-12 bg-gradient-to-r from-stellar-500 to-primary-600 rounded-lg flex items-center justify-center">
          <Users className="w-6 h-6 text-white" />
        </div>
       