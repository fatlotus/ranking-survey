import json
import random

# urls = [
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1444.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1422.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1424.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1428.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1452.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1449.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1442.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1454.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1453.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1433.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1445.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1451.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1435.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1427.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1423.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1441.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1443.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1434.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1431.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1432.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1426.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1425.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1450.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1429.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1430.jpg",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1446.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1447.png",
#     "http://flyonovereq.com/wp-content/uploads/2015/10/IMG_1448.png"
# ]

urls = [
    "http://i.giphy.com/HhvUpQhBWMtwc.gif",
    "http://i.giphy.com/1rBCI5HKJPd0k.gif",
    "http://i.giphy.com/undefined.gif",
    "http://i.giphy.com/ZqcEzvXpivVwk.gif",
    "http://i.giphy.com/l2JJD7JIC2hZrr2Hm.gif",
    "http://i.giphy.com/Otu4kaNHJgrqo.gif",
    "http://i.giphy.com/z3MVNiOylhon6.gif",
    "http://i.giphy.com/ah7KwjMNJlhtK.gif",
    "http://i.giphy.com/xTiTnFjIvn3fqxGwA8.gif",
    "http://i.giphy.com/gZ7UAPE4SnjJC.gif",
    "http://i.giphy.com/W8hR8uEvHjane.gif",
    "http://i.giphy.com/3HGrwg1lDWzxS.gif",
    "http://i.giphy.com/gPfQ1nFixU4pi.gif",
    "http://i.giphy.com/u7JT9EcX9RAZy.gif",
    "http://i.giphy.com/3xz2BO2QyUfHRmV9DO.gif",
    "http://i.giphy.com/l3V0xVRGPKoSidlUA.gif",
    "http://i.giphy.com/7Oo6pOoNDYXni.gif",
    "http://i.giphy.com/fW9RtiTX4LxhS.gif",
    "http://i.giphy.com/119zR8FNY0h6lq.gif",
    "http://i.giphy.com/N18WLp1lduj2U.gif",
    "http://i.giphy.com/l0NhWOpqU20G6GFmU.gif",
    "http://i.giphy.com/j7c8Gcxet30Iw.gif",
    "http://i.giphy.com/Kbbt8JCKCoYWA.gif",
    "http://i.giphy.com/kTQATxy3QZveE.gif",
    "http://i.giphy.com/RM3qdSRUul3ag.gif",
    "http://i.giphy.com/3o6ozvFZGdqpIakCQg.gif",
    "http://i.giphy.com/orsKKFqWKtsVa.gif",
    "http://i.giphy.com/eaEl2EIDmHE0U.gif",
    "http://i.giphy.com/UmNHGGGcvCHvi.gif",
    "http://i.giphy.com/iQidCLEchWI4o.gif",
    "http://i.giphy.com/qwqD0ycLPxyG4.gif",
    "http://i.giphy.com/26tn5FsLe5oK6VKVy.gif",
    "http://i.giphy.com/EcjaT6tBBB196.gif",
    "http://i.giphy.com/3o6gbcoRtMf6LLwWBi.gif",
    "http://i.giphy.com/O6EeEC6TC6sw.gif",
    "http://i.giphy.com/3o6ozEfvj2hIzyDkvC.gif",
    "http://i.giphy.com/qt0Frx4shI35e.gif",
    "http://i.giphy.com/4IdRoU2fx0kRG.gif",
    "http://i.giphy.com/t1cDOoX8xgWEo.gif",
    "http://i.giphy.com/3o85g2nHOWjdhhefHq.gif",
    "http://i.giphy.com/2qlPyvkTvKDSM.gif",
    "http://i.giphy.com/dgPkAhGr3bBks.gif",
    "http://i.giphy.com/mDzplQueBILHG.gif",
    "http://i.giphy.com/OWRejDOkTgX6M.gif",
    "http://i.giphy.com/y9EayOkRK1WN2.gif",
    "http://i.giphy.com/xT1XGYhZZlnrgqETDi.gif",
    "http://i.giphy.com/BnSfEHS0ix3Ve.gif",
    "http://i.giphy.com/Fd5rhgeRn2LsY.gif",
    "http://i.giphy.com/xneZmsVUJSYb6.gif",
    "http://i.giphy.com/dQX6YhEcJuEes.gif",
    "http://i.giphy.com/3o6ozk957G2FE4mMCY.gif",
    "http://i.giphy.com/EyrcMdlo9IrM4.gif",
    "http://i.giphy.com/3o6gaSgLGfCWdCrlCg.gif",
    "http://i.giphy.com/C6vNLxPgSIheo.gif",
    "http://i.giphy.com/VARZEXQKRvhDO.gif",
    "http://i.giphy.com/3o6ozpaDcskZKX2pLq.gif",
    "http://i.giphy.com/XDqraU5tKYqJi.gif",
    "http://i.giphy.com/VmNVXIMH3bCsE.gif",
    "http://i.giphy.com/xT1XGzd8Up0nZBS0Gk.gif",
    "http://i.giphy.com/8rAulVj0FomaI.gif",
    "http://i.giphy.com/RMyjqBYU20I6Y.gif",
    "http://i.giphy.com/Fw4BSHB7UZ560.gif",
    "http://i.giphy.com/r8Kn8NaNM3xba.gif",
    "http://i.giphy.com/3oEduPs99YeI3xBD3y.gif",
    "http://i.giphy.com/KfyJtHHe0uZiw.gif",
    "http://i.giphy.com/hTgeSxaiyvYK4.gif",
    "http://i.giphy.com/okfvUCpgArv3y.gif",
    "http://i.giphy.com/LNWqOKuh6SJKo.gif",
    "http://i.giphy.com/ALCI3eTii7qOk.gif",
    "http://i.giphy.com/C7G76lpbpRvC8.gif",
    "http://i.giphy.com/HlJ6JD8XpWfkc.gif",
    "http://i.giphy.com/gxjiNJ2XPA8py.gif",
    "http://i.giphy.com/INtH4d27qaKIM.gif",
    "http://i.giphy.com/TRMFOJbupo4aQ.gif",
    "http://i.giphy.com/9Y6n9TR7U07ew.gif",
    "http://i.giphy.com/ALneyoThCd8lO.gif",
    "http://i.giphy.com/l4Ki5elhhHcF4jvHi.gif",
    "http://i.giphy.com/3osxYjaU51nj1drij6.gif",
    "http://i.giphy.com/iwJMmqOiqzss0.gif",
    "http://i.giphy.com/3osxYx5aUbpiQj0cWk.gif",
    "http://i.giphy.com/xT1XH35xaGNfnYQUYo.gif",
    "http://i.giphy.com/JXvN0AvZTpD2w.gif",
    "http://i.giphy.com/T0S33Gx1DcCdi.gif",
    "http://i.giphy.com/8H6ZAoM8wXmes.gif",
    "http://i.giphy.com/gzBv9sBO1OFfW.gif",
    "http://i.giphy.com/pHx6lgq3fiRb2.gif",
    "http://i.giphy.com/K791JhEoRS6QM.gif",
    "http://i.giphy.com/2HkiO2q64Dg40.gif",
    "http://i.giphy.com/n2jGtmTGMEvpS.gif",
    "http://i.giphy.com/3D6sVP9MCEn7i.gif",
    "http://i.giphy.com/3o85xr3K9IXEWAb4Eo.gif",
    "http://i.giphy.com/rYYYpc8pGX7tS.gif",
    "http://i.giphy.com/nzKbmkN3i8RTq.gif",
    "http://i.giphy.com/Z4XJsH9mUZUY.gif",
    "http://i.giphy.com/xT1XGPBvbwPBlphBSM.gif",
    "http://i.giphy.com/5Xqq1IfqFg8X6.gif",
    "http://i.giphy.com/UMvAwWfZo8zW8.gif",
    "http://i.giphy.com/xT0BKK3QJwa4csEzD2.gif",
    "http://i.giphy.com/3o6gaZiWBQbZw7MBUI.gif",
    "http://i.giphy.com/3o6gb2Oy0bdNpdxNfO.gif",
    "http://i.giphy.com/ffCgISyzuly0w.gif",
    "http://i.giphy.com/26AHLMsoGBRExNI08.gif",
    "http://i.giphy.com/8gi3P2vS0FF6M.gif",
    "http://i.giphy.com/3o6ozCAEzM2jgeMyqc.gif",
    "http://i.giphy.com/3o6ozC0tdjoudmcXWE.gif",
    "http://i.giphy.com/xT1XGHWeSM81vT9lWo.gif",
    "http://i.giphy.com/26AHGkrF8DnmyzxAY.gif",
    "http://i.giphy.com/l0DEKUFABhqS6PpiU.gif",
    "http://i.giphy.com/xT1XGB8VnRGgU8E86I.gif",
    "http://i.giphy.com/xTiTnsQvlHAQkOafio.gif",
    "http://i.giphy.com/hkut0Q7K6GPle.gif",
    "http://i.giphy.com/5fBH6zgbXfijdW2ypkQ.gif",
    "http://i.giphy.com/3o6ozkS0d2wOcl8naE.gif",
    "http://i.giphy.com/Gt8iOfgzOf0WI.gif",
    "http://i.giphy.com/3osxYp2s5tsJWTnpwA.gif",
    "http://i.giphy.com/3osxYmDloqKIFWVCMg.gif",
    "http://i.giphy.com/66XPTwbYSuIaQ.gif",
    "http://i.giphy.com/3o6ozkmZHCEcuqczF6.gif",
    "http://i.giphy.com/uHgH0FVgnhCw0.gif",
    "http://i.giphy.com/qgVuPuTMTK0OA.gif",
    "http://i.giphy.com/zdFNkQ6vGJADC.gif",
    "http://i.giphy.com/FELYJX2AVLuc8.gif",
    "http://i.giphy.com/13wkIMZF3gdfwc.gif",
    "http://i.giphy.com/26u6bDoZ7ReueAqeQ.gif",
    "http://i.giphy.com/3oJpxQGBWmoWGlscxy.gif",
    "http://i.giphy.com/gQydVHW41aQow.gif",
    "http://i.giphy.com/rE04CIrzIrSN2.gif",
    "http://i.giphy.com/ZbIWtQFH6KJeo.gif",
    "http://i.giphy.com/NNrrgtWVB3WZW.gif",
    "http://i.giphy.com/107wctifgEsaqI.gif",
    "http://i.giphy.com/qdDxwZxamaVHy.gif",
    "http://i.giphy.com/l3V0rsC3FUS106CKk.gif",
    "http://i.giphy.com/26AHJtympt9jypus8.gif",
    "http://i.giphy.com/3o6ozkEfP7bhvIUi3u.gif",
    "http://i.giphy.com/xT1XGCkmFVqXHbmb60.gif",
    "http://i.giphy.com/xT1XH08vQSb0YJibDy.gif",
    "http://i.giphy.com/xTiTnpCQilvKMrIXrW.gif",
    "http://i.giphy.com/ypbVCkpAVA6fm.gif",
    "http://i.giphy.com/zxMq6nR6lFVFC.gif",
    "http://i.giphy.com/6oMhxODhXtyW4.gif",
    "http://i.giphy.com/xTiTnHz4LKeY7zNDvq.gif",
    "http://i.giphy.com/xThuW3Xz08v5kMfBUk.gif",
    "http://i.giphy.com/xThuW0IW1lhMlcDTEc.gif",
    "http://i.giphy.com/26AHOQHEXuHba91hC.gif",
    "http://i.giphy.com/3osxYxSkryiJW8GYZq.gif",
    "http://i.giphy.com/YbXtbKoi2ZUOc.gif",
    "http://i.giphy.com/LrN9NbJNp9SWQ.gif",
    "http://i.giphy.com/4kA6QUa0cd60g.gif",
    "http://i.giphy.com/Ifgh6AsnGxbgY.gif",
    "http://i.giphy.com/xT1XGD02MRIfS829qw.gif",
    "http://i.giphy.com/3o7WTGHSxNEvAvIOqs.gif",
    "http://i.giphy.com/VIfi3NzeoDVT2.gif",
    "http://i.giphy.com/gCRTm06QZZEek.gif",
    "http://i.giphy.com/aB2xXDE6V94Zy.gif",
    "http://i.giphy.com/QBC5foQmcOkdq.gif",
    "http://i.giphy.com/xT1XGQlwXWcRZOdBcI.gif",
    "http://i.giphy.com/3oEduM6VI2dbST7F2E.gif",
    "http://i.giphy.com/xT1XGG2IJIzwlmsAkU.gif",
    "http://i.giphy.com/xT1XGQve0LxCblDLr2.gif",
    "http://i.giphy.com/51msWHqr8drws.gif",
    "http://i.giphy.com/W8krmZSDxPIfm.gif",
    "http://i.giphy.com/pF7vnUwzDZPFu.gif",
    "http://i.giphy.com/kHGNWt9O0L8Zi.gif",
    "http://i.giphy.com/7sKtfIAbaATx6.gif",
    "http://i.giphy.com/l3V0dclK0lcUw2ygE.gif",
    "http://i.giphy.com/l3V0BCbe9YgCZN8hW.gif",
    "http://i.giphy.com/12xiOA46vEYGl2.gif",
    "http://i.giphy.com/3osxYvRjG7SLCAdILm.gif",
    "http://i.giphy.com/LEdz8xl9uFxKw.gif",
    "http://i.giphy.com/l3V0i1UL660cLFPFu.gif",
    "http://i.giphy.com/3o6ozifkiTecWsJOFO.gif",
    "http://i.giphy.com/26AHMkSi8F8wgsHzW.gif",
    "http://i.giphy.com/CovFciJgWyxUs.gif"]

markup = ['<img src="{}" width="95%"/>'.format(url) for url in urls]

with open("questions.json", "w") as fp:
    for j in xrange(100):
        for subject in [2, 3, 5, "cmp"]:
            for continuous in [False, True]:
                survey = "experiment/{}/{}{}".format(
                    j, subject, "-cont" if continuous else "")
                comparisons, rating = 20, 80
                if subject == "cmp":
                    comparisons, rating = 100, 0

                for i in xrange(comparisons):
                    a = random.choice(markup)
                    b = random.choice(markup)

                    json.dump({
                        "survey": survey,
                        "choices": [a, b],
                        "precision": 2,
                        "exclusive": True,
                    }, fp)
                    fp.write("\n")

                for i in xrange(rating):
                    a = random.choice(markup)

                    json.dump({
                        "survey": survey,
                        "choices": [a],
                        "precision": subject * 100 if continuous else subject,
                        "exclusive": True,
                    }, fp)
                    fp.write("\n")

    # for i in xrange(100):
    #     a = random.choice(markup)
    #     b = random.choice(markup)
    #     c = random.choice(markup)
    #
    #     json.dump({
    #         "survey": "ranking",
    #         "choices": [a, b, c],
    #         "precision": 3,
    #         "exclusive": True,
    #     }, fp)
    #     fp.write("\n")

    # for i in xrange(100):
    #     a = random.choice(markup)
    #
    #     json.dump({
    #         "survey": "rating2",
    #         "choices": [a],
    #         "precision": 2,
    #         "exclusive": False,
    #     }, fp)
    #     fp.write("\n")

    # for i in xrange(100):
    #     a = random.choice(markup)
    #
    #     json.dump({
    #         "survey": "rating3",
    #         "choices": [a],
    #         "precision": 3,
    #         "exclusive": False,
    #     }, fp)
    #     fp.write("\n")

    # for i in xrange(100):
    #     a = random.choice(markup)
    #
    #     json.dump({
    #         "survey": "rating5",
    #         "choices": [a],
    #         "precision": 5,
    #         "exclusive": False,
    #     }, fp)
    #     fp.write("\n")
