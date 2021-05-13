# FDALang
интерпретатор для игры FDA

Интерпретатор для простого, но строгого языка FDALang, имеющим следующие особенности:
* 1 стейтмент на одну строку (исключение блочные стейтменты типа if/switch/for). Стейтмент это выражение, которое не возвращает результат
* исходя из пункта выше, стейтменты не нужно завершать символом `;`
* язык со строгой типизацией, но без объявления переменных - тип определяется при инициализации, и не может быть впоследствии изменен
* нельзя проводить операции над разными типами, даже если это float и int - будет ошибка. нужно использовать приведение типов типа `a = 3 + int(4.5)`
* функции всегда задаются как переменные для простоты синтаксиса
* Go/Cи-подобный синтаксис, но без указателей
* Возможность указывать тип с пустым значением, это типа как null, только типизированный
* примеры простых программ:
```
sum = fn(int x, int y) int {
   return x + y
}
a = sum(2, 5)
c = 10
if c > 8 {
    print(true)
} else {
    print(false)
}
struct point {
   float x
   float y
}
p = point{x = 1., y = 2.}
px = p.x
```

пример программы для игры, базовые действия:
```
commands.move = 1.
commands.rotate = 0.5
commands.cannon.rotate = -0.8
commands.cannon.shoot = 0.1
```

пример программы для реальной игры:
```
xelon = getFirstTarget(1)
if empty(xelon) {
   return 1
}
angleObj = angle(mech.x, mech.y, obj.x, obj.y)
angleMech = mech.angle
angleTo = angleObj - angleMech
if angleTo < -PI {
   angleTo = 2. * PI + angleTo
}
if angleTo > PI {
   angleTo = angleTo - 2. * PI
}

switch angleTo {
case > 1.:
   commands.rotate = 1.
case < -1.:
   commands.rotate = -1.
default:
   commands.rotate = angleTo
}

dist = distance(mech.x, mech.y, obj.x, obj.y)
if obj.type == 3 {
   commands.move = 1.
   return 1
}
if dist > 200. {
   commands.move = distance / 1000.
   if commands.move > 1. {
      commands.move = 1.
   }
}
toShoot = commands.rotate * commands.rotate * dist
if toShoot < 70. {
   commands.cannon.shoot = 0.1
   return 1
}
```

пример чуть посложнее:
```
xelon = getFirstTarget(1)
if empty(xelon) {
   xelon = nearestByType(mech, objects, ObjectTypes:xelon)
   if empty(xelon) {
      return 1
   }
   addTarget(xelon, 1)
}
angleTo = angleToRotate(mech.angle, mech.x, mech.y, xelon.x, xelon.y)
commands.rotate = keepBounds(angleTo, 1.)

commands.move = 1. - absFloat(commands.rotate)

obj = getFirstTarget(2)
if empty(obj) {
   obj = nearestByType(mech, objects, ObjectTypes:spore)
   if empty(obj) {
      return 1
   }
   addTarget(obj, 2)
}

angleSum = mech.angle + mech.cAngle
cAngleTo = angleToRotate(angleSum, mech.x, mech.y, obj.x, obj.y)

if cAngleTo * angleTo < 0. {
   cAngleTo = cAngleTo - angleTo
}
commands.cannon.rotate = keepBounds(cAngleTo, 1.)

dist = distance(mech.x, mech.y, obj.x, obj.y)
toShoot = cAngleTo * cAngleTo * dist
if toShoot < 40. {
   commands.cannon.shoot = 0.1
}
```

# TODO
* control flow - for
* Тип string
* Поддержка пакетов
* Контроль глубины стэка вызовов
* Бенчмарки - трэкинг производительности интерпретатора
* стэктрейс при ошибках
* Импорты
