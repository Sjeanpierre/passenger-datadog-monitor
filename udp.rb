require 'socket'
BasicSocket.do_not_reverse_lookup = true
# Create socket and bind to address
client = UDPSocket.new
client.bind(nil, 8125)
loop do
  data, addr = client.recvfrom(1024) # if this number is too low it will drop the larger packets and never give them to you
  puts "From addr: '%s', msg: '%s'" % [addr.join(','), data]
end
