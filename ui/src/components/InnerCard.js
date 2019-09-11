import React from 'react'
import './../App.css'
import { Card, majorScale } from 'evergreen-ui'

const InnerCard = ({ children }) => (
  <Card
    display="flex"
    flexDirection="column"
    elevation={1}
    background="white"
    marginTop={majorScale(2)}
  >
    {children}
  </Card>
)

export default InnerCard