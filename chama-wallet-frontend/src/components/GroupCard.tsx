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

 