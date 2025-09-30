import React, { useState } from 'react'
import { Users, Crown, Vote } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { groupApi } from '../services/api'
import type { Group, Member, User } from '../types'
import toast from 'react-hot-toast'

interface AdminNominationProps {
  group: Group
  currentUser: User
}

const AdminNomination: React.FC<AdminNominationProps> = ({ group, currentUser }) => {
  const [showNominationModal, setShowNominationModal] = useState(false)
  const [selectedMember, setSelectedMember] = useState('')
  const queryClient = useQueryClient()

  const nominateAdminMutation = useMutation({
    mutationFn: (data: { nominee_id: string }) => 
      groupApi.nominateAdmin(group.ID, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Admin nomination submitted!')
      setShowNominationModal(false)
      setSelectedMember('')
    },
  