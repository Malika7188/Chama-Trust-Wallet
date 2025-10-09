import React, { useState } from 'react'
import { DollarSign, Check, X, Clock, Users } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { groupApi, payoutApi } from '../services/api'
import type { Group, User, PayoutRequest } from '../types'
import toast from 'react-hot-toast'

interface PayoutManagementProps {
  group: Group
  currentUser: User
}

