interface LoadingSpinnerProps {
  
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

const LoadingSpinner = ({ size = 'md', className = '' }: LoadingSpinnerProps) => {
  