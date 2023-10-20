import machine
import utime
import time
import network
import socket
import _thread
import uasyncio
from machine import Pin

##### Pin Definition
servo = machine.PWM(machine.Pin(15))
button = machine.Pin(14,machine.Pin.IN)
led = Pin("LED", machine.Pin.OUT)

##### Global Valiable Definition
HOST='192.168.11.200'

open = 0
close = 180
openFlag = True # True is Open
operateFlag = False # True is in Operating


##### Servo & Button Operator
def interval_mapping(x, in_min, in_max, out_min, out_max):
    return (x - in_min) * (out_max - out_min) / (in_max - in_min) + out_min

def servo_write(pin, angle):
    pin.freq(50)
    pulse_width = interval_mapping(angle, 0, 180, 0.5, 2.5)
    duty = int(interval_mapping(pulse_width, 0, 20, 0, 65535))
    pin.duty_u16(duty)

def close_key():
    global operateFlag
    global openFlag
    global led
    
    if operateFlag:
        return "another",False
    
    if not openFlag:
        servo_write(servo,close)
        utime.sleep(3)
        return "already",False
    
    operateFlag = True
    servo_write(servo,close)
    openFlag=False
    led.value(0)
    utime.sleep(3)
    servo.deinit()
    operateFlag = False
    return "complete",False
    
def open_key():
    global operateFlag
    global openFlag
    global led
    
    if operateFlag:
        return "another",False
    
    if openFlag:
        servo_write(servo,open)
        utime.sleep(3)
        return "already",True
    
    operateFlag = True
    servo_write(servo,open)
    openFlag = True
    led.value(1)
    utime.sleep(3)
    servo.deinit()
    operateFlag = False
    return "complete",True

def init_servo():
    global openFlag
    global operateFlag
    global led
    operateFlag = True
    openFlag = True
    
    print('init')
    servo_write(servo,open)
    led.value(1)
    utime.sleep_ms(3000)
    servo.deinit()
    operateFlag = False

def button_operator():
    if button.value() == 1:
        print("operate")
        if operateFlag:
            return
        
        if openFlag:
            close_key()
        else:
            open_key()
##### end Servo & Button Operator



##### start Server Block
# Listen for connections, serve client
def server_operator(s):
    try:
        print('socket ready!')
        cl, addr = s.accept()
        
        print('client connected from', addr)
        print('handle Request')
        request = cl.recv(1024)
        print("request:")
        request = str(request)
        print(request)

        opStatus = "unknown"
        keyStatus = "False"
        
        if '/open' in request:
            print('open request')
            opStatus,keyStatus = open_key()
        elif '/close' in request:
            print('close request')
            opStatus,keyStatus = close_key()
        elif '/check' in request:
            print('check request')
            if operateFlag :
                opStatus = "another"
            else :
                opStatus = "already"
            keyStatus = openFlag

        else:
            print('!!!!! Unknown Request !!!!!')
        
        body = '{"keyStatus":"'+str(keyStatus)+'","opStatus":"'+opStatus+'"}'

        response = 'HTTP/1.0 200 OK\r\nContent-type: text/html\r\n\r\n'+body
        cl.send(bytes(response,'UTF-8'))
        cl.close()
        
        
    except OSError as e:
        pass

def main() :
    #自宅Wi-FiのSSIDとパスワードを入力
    ssid = 'BarleyTea'
    password = 'kuma1221'

    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)
    wlan.connect(ssid, password)
    wlan.ifconfig(('192.168.11.200', '255.255.255.0', '192.168.0.1', '8.8.8.8'))

    # Wait for connect or fail
    max_wait = 10
    while max_wait > 0:
        if wlan.status() < 0 or wlan.status() >= 3:
            break
        max_wait -= 1
        print('waiting for connection...')
        time.sleep(1)
        
    # Handle connection error
    if wlan.status() != 3:
        raise RuntimeError('network connection failed')
    else:
        print('Connected')
        status = wlan.ifconfig()
        print( 'ip = ' + status[0] )
        
        
    # Open socket
    addr = socket.getaddrinfo('0.0.0.0', 80)[0][-1]
    s = socket.socket()
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1) 
    s.bind(addr)
    s.settimeout(1)
    s.listen(1)
    print('listening on', addr)
    
    # Start Button Thread
    #_thread.start_new_thread(button_operator,())
    init_servo()
    while True:
        server_operator(s)
        button_operator()

main()

