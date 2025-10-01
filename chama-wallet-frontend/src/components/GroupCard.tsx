import { Link } from 'react-router-dom';
import { Users } from 'lucide-react';
import { useGroupBalance } from '../hooks/useGroups';
import type { Group } from '../types';

interface GroupCardProps {
  group: Group;
}

