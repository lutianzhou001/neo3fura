using Neo;
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Native;
using Neo.SmartContract.Framework.Services;
using System;
using System.Numerics;

namespace Neoverse
{
    public partial class Neoverse
    {
        private static readonly StorageMap OwnerMap = new(Storage.CurrentContext, 0x16);

        private static bool IsOwner() => Runtime.CheckWitness(GetOwner());

        public static bool Verify() => IsOwner();

        public static UInt160 SetOwner(UInt160 newOwner)
        {
            if (!IsOwner()) throw new Exception("No authorization.");
            if (!newOwner.IsValid) throw new Exception("Neoverse::SetOwner: UInt160 is invalid.");

            OwnerMap.Put("owner", newOwner);
            return GetOwner();
        }

        public static void _deploy(object _, bool update)
        {
            if (update) return;
            IndexStorage.Initial();
            FragmentStorage.Initial();
            CounterStorage.Initial();
        }

        public static void Update(ByteString nefFile, string manifest)
        {
            if (!IsOwner()) throw new Exception("No authorization.");
            ContractManagement.Update(nefFile, manifest, null);
        }

        public static void Destroy()
        {
            if (!IsOwner()) throw new Exception("No authorization.");
            ContractManagement.Destroy();
        }
        public static bool Pause()
        {
            if (!IsOwner()) throw new Exception("No authorization.");
            StateStorage.Pause();
            return true;
        }

        public static bool Resume()
        {
            if (!IsOwner()) throw new Exception("No authorization.");
            StateStorage.Resume();
            return true;
        }

        public static bool Airdrop(UInt160 to, BigInteger amount)
        {
            if (!IsOwner()) throw new Exception("No authorization.");
            if (!to.IsValid) throw new Exception("Neoverse::Airdrop: UInt160 is invalid.");

            for (int i = 0; i < amount; i++)
            {
                var number = IndexStorage.NextIndex(0);
                if (number > FragmentStorage.MaxCount * 9) throw new Exception("Neoverse::Airdrop: Sold out.");
                var blindBox = TokenState.CreateBlindBox(to, number);
                Mint(blindBox.Name, blindBox);
            }
            return true;
        }
    }
}
