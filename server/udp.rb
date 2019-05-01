require 'socket'

BasicSocket.do_not_reverse_lookup = true

# Create socket and bind to address
client = UDPSocket.new
client.bind('0.0.0.0', 8125)

loop do
  # if this number is too low it will drop the larger packets and never give them to you
  data, addr = client.recvfrom(1024)
  puts "From addr: '%s', msg: '%s'" % [addr.join(','), data]
end
