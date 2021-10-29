using Neo;
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services;
using System;
using System.Numerics;

namespace Neoverse
{
    public class TokenState : Nep11TokenState
    {
        public byte FirstType;      //0盲盒 1碎片 2卡牌
        public byte SecondType;     //盲盒：0，碎片A-I：0-8，卡牌N、E、O：0-2
        public BigInteger Number;   //盲盒、碎片、卡牌的编号
        public string Image;        //图片URL，中心化存储 e.g. https://neo.org/BlindBox.png
        public string Video;        //视频URL，中心化存储 e.g. https://neo.org/BlindBox.mp4

        /// <summary>
        /// 创建盲盒
        /// </summary>
        /// <param name="owner">拥有者</param>
        /// <param name="randomIndex">下一个区块的Index</param>
        /// <param name="number">盲盒编号</param>
        public static TokenState CreateBlindBox(UInt160 owner, BigInteger number) => new(owner, 0, 0, number);

        /// <summary>
        /// 创建碎片
        /// </summary>
        /// <param name="owner">拥有者</param>
        /// <param name="fragmentType">碎片A-I：0-8</param>
        /// <param name="number">碎片编号</param>
        public static TokenState CreateFragment(UInt160 owner, byte fragmentType, BigInteger number) => new(owner, 1, fragmentType, number);

        /// <summary>
        /// 创建卡牌
        /// </summary>
        /// <param name="owner">拥有者</param>
        /// <param name="cardType">卡牌N、E、O：0-2</param>
        /// <param name="number">卡牌编号</param>
        public static TokenState CreateCard(UInt160 owner, byte cardType, BigInteger number) => new(owner, 2, cardType, number);

        private TokenState(UInt160 owner, byte firstType, byte secondType, BigInteger number)
        {
            Owner = owner;
            FirstType = firstType;
            SecondType = secondType;
            Number = number;
            
            switch (firstType)
            {
                case 0: Name = "Blind Box #" + number;
                    Image = "https://neo.org/BlindBox.png";
                    Video = "https://neo.org/BlindBox.mp4";
                    break;
                case 1:
                    var typeName = (ByteString)((BigInteger)'A' + secondType);
                    Name = "Fragment " + typeName + " #" + number;
                    Image = "https://neo.org/Fragment" + typeName + ".png";
                    Video = "https://neo.org/Fragment" + typeName + ".mp4";
                    break;
                case 2:
                    var index = (byte)(number % 3);
                    var imageName = new string[0];
                    switch (secondType)
                    {
                        case 0:
                            imageName = new string[] { "Xira", "Aeras", "Nero" };
                            Name = "N #" + number + " " + imageName[index]; break;
                            // N #1 Aeras
                            // N #2 Nero
                            // N #3 Xira
                        case 1:
                            imageName = new string[] { "Noiz", "Zion", "Core" };
                            Name = "E #" + number + " " + imageName[index]; break;
                            // E #101 Core
                            // E #102 Noiz
                            // E #103 Zion
                        case 2:
                            imageName = new string[] { "Interoperability", "Composability", "Scalability" };
                            Name = "O #" + number + " " + imageName[index]; break;
                            // O #301 Interoperability
                            // O #302 Composability
                            // O #303 Scalability
                    }
                    Image = "https://neo.org/" + imageName[index] + ".png";
                    break;
            }
        }

        public void CheckAdmin()
        {
            if (Runtime.CheckWitness(Owner)) return;
            throw new InvalidOperationException("No authorization.");
        }
    }
}
